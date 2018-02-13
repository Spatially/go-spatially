package spatiallydb

import (
	"bytes"
	"fmt"
	"reflect"

	geojson "github.com/paulmach/go.geojson"
)

// WKTToGeometry converts a given Well Known Text shape into a geojson geometry
func WKTToGeometry(wkt string) (g *geojson.Geometry, err error) {
	return parseWKT([]byte(wkt))
}

func parseWKT(data []byte) (g *geojson.Geometry, err error) {
	s := &scanner{raw: data}
	return s.scanGeom()
}

type scanner struct {
	raw []byte
	i   int
}

func (s *scanner) scanGeom() (g *geojson.Geometry, err error) {
	ident, err := s.scanIndent()
	if err != nil {
		return g, err
	}
	switch ident {
	case "POINT", "MULTIPOINT", "LINESTRING":
		var ps [][]float64
		ps, err = s.scanPoints(ident == "MULTIPOINT")
		if err != nil {
			break
		}
		switch ident {
		case "POINT":
			if len(ps) != 1 {
				return nil, fmt.Errorf("expected 1 got %d points", len(ps))
			}
			g = geojson.NewPointGeometry(ps[0])
		case "MULTIPOINT":
			g = geojson.NewMultiPointGeometry(ps...)
		case "LINESTRING":
			g = geojson.NewLineStringGeometry(ps)
		}
	case "POLYGON", "MULTILINESTRING":
		isPolygon := true
		if ident != "POLYGON" {
			isPolygon = false
		}
		var polygon [][][]float64
		polygon, err = s.scanMultiLineString(isPolygon)
		if err != nil {
			return nil, err
		}
		if isPolygon {
			g = geojson.NewPolygonGeometry(polygon)
		} else {
			g = geojson.NewLineStringGeometry(polygon[0])
		}
	case "MULTIPOLYGON":
		var multipolygon [][][][]float64
		multipolygon, err = s.scanMultiPolygon()
		if err != nil {
			return nil, err
		}
		g = geojson.NewMultiPolygonGeometry(multipolygon...)
	default:
		err = fmt.Errorf("unknown or unimplemented geometry '%s'", ident)
	}
	if err != nil {
		return nil, err
	}
	return
}

func (s *scanner) scanIndent() (string, error) {
	s.skipWs()
	var b byte
	i, start := 0, -1
	for i, b = range s.raw[s.i:] {
		if b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z' {
			if start < 0 {
				start = i
			}
			continue
		}
		if start < 0 {
			return "", fmt.Errorf("no ident '%v'", b)
		}
		break
	}
	str := string(s.raw[s.i+start : i-start])
	s.i += i
	return str, nil
}

func (s *scanner) skipWs() {
	for i, b := range s.raw[s.i:] {
		if b == ' ' || b == '\n' || b == '\t' || b == '\r' || b == 65 {
			continue
		}
		s.i += i
		return
	}
}

func (s *scanner) scanPoints(multi bool) ([][]float64, error) {
	err := s.scanStart()
	if err != nil {
		return nil, err
	}
	var ps [][]float64
	var p []float64
	var comma bool
	if multi {
		err = s.scanStart()
		multi = err == nil
	}
	for {
		p, comma, err = s.scanPoint()
		if err != nil {
			return nil, err
		}
		if len(p) < 2 {
			return nil, fmt.Errorf("point must be at least 2d. got %d elements", len(p))
		}
		if len(p) < 2 || len(p) > 4 {
			return nil, fmt.Errorf("point can be at most 4d. got %d elements", len(p))
		}
		ps = append(ps, p)
		if comma {
			if multi {
				return nil, fmt.Errorf("expect ')' got ','")
			}
			continue
		}
		if multi {
			comma, err = s.scanContinue()
			if err != nil {
				return nil, err
			}
			if comma {
				err = s.scanStart()
				if err != nil {
					return nil, err
				}
				continue
			}
		}
		return ps, nil
	}
}

func (s *scanner) scanStart() error {
	s.skipWs()
	start := s.raw[s.i] == '('
	if !start {
		return fmt.Errorf("expect '(' got '%v'", s.raw[s.i])
	}
	s.i++
	return nil
}

func (s *scanner) scanPoint() (p []float64, comma bool, err error) {
	var pc []float64
	s.skipWs()
	r := bytes.NewReader(s.raw[s.i:])
	var f float64
	for {
		_, err := fmt.Fscan(r, &f)
		if err != nil {
			return nil, false, err
		}
		pc = append(pc, f)
		s.i = len(s.raw) - r.Len()
		s.skipWs()
		b := s.raw[s.i]
		if comma = b == ','; comma || b == ')' {
			s.i++
			break
		}
	}
	return pc, comma, nil
}

func (s *scanner) scanContinue() (bool, error) {
	s.skipWs()
	comma := s.raw[s.i] == ','
	if !comma && s.raw[s.i] != ')' {
		return comma, fmt.Errorf("expect ',' or ')' got '%v'", s.raw[s.i])
	}
	s.i++
	return comma, nil
}

func (s *scanner) scanMultiLineString(isPolygon bool) ([][][]float64, error) {
	err := s.scanStart()
	if err != nil {
		return nil, err
	}
	var mls [][][]float64
	var ps [][]float64
	var comma bool
	for {
		ps, err = s.scanPoints(false)
		if err != nil {
			return nil, err
		}
		if isPolygon {
			if len(ps) < 4 {
				return nil, fmt.Errorf("a polygon must have at least 4 points, got %d", len(ps))
			}
			if !reflect.DeepEqual(ps[0], ps[len(ps)-1]) {
				return nil, fmt.Errorf("a polygon must be closed")
			}
		}
		mls = append(mls, ps)
		comma, err = s.scanContinue()
		if err != nil {
			return nil, err
		}
		if comma {
			continue
		}
		return mls, nil
	}
}

func (s *scanner) scanMultiPolygon() ([][][][]float64, error) {
	err := s.scanStart()
	if err != nil {
		return nil, err
	}
	var multi [][][][]float64
	var poly [][][]float64
	var comma bool
	for {
		poly, err = s.scanMultiLineString(true)
		if err != nil {
			return nil, err
		}
		multi = append(multi, poly)
		comma, err = s.scanContinue()
		if err != nil {
			return nil, err
		}
		if comma {
			continue
		}
		return multi, nil
	}
}
