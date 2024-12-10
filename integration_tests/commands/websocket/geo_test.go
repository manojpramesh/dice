package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeoAdd(t *testing.T) {
	exec := NewWebsocketCommandExecutor()
	conn := exec.ConnectToServer()

	testCases := []struct {
		name   string
		cmds   []string
		expect []interface{}
	}{
		{
			name:   "GeoAdd With Wrong Number of Arguments",
			cmds:   []string{"GEOADD points 1 2"},
			expect: []interface{}{"ERR wrong number of arguments for 'geoadd' command"},
		},
		{
			name:   "GeoAdd With Adding New Member And Updating it",
			cmds:   []string{"GEOADD points 1.21 1.44 NJ", "GEOADD points 1.22 1.54 NJ"},
			expect: []interface{}{float64(1), float64(0)},
		},
		{
			name:   "GeoAdd With Adding New Member And Updating it with NX",
			cmds:   []string{"GEOADD points NX 1.21 1.44 MD", "GEOADD points 1.22 1.54 MD"},
			expect: []interface{}{float64(1), float64(0)},
		},
		{
			name:   "GEOADD with both NX and XX options",
			cmds:   []string{"GEOADD points NX XX 1.21 1.44 DEL"},
			expect: []interface{}{"ERR XX and NX options at the same time are not compatible"},
		},
		{
			name:   "GEOADD invalid longitude",
			cmds:   []string{"GEOADD points 181.0 1.44 MD"},
			expect: []interface{}{"ERR invalid longitude"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i, cmd := range tc.cmds {
				result, err := exec.FireCommandAndReadResponse(conn, cmd)
				assert.Nil(t, err)
				assert.Equal(t, tc.expect[i], result, "Value mismatch for cmd %s", cmd)
			}
		})
	}
}

func TestGeoDist(t *testing.T) {
	exec := NewWebsocketCommandExecutor()
	conn := exec.ConnectToServer()
	defer conn.Close()

	testCases := []struct {
		name   string
		cmds   []string
		expect []interface{}
	}{
		{
			name: "GEODIST b/w existing points",
			cmds: []string{
				"GEOADD points 13.361389 38.115556 Palermo",
				"GEOADD points 15.087269 37.502669 Catania",
				"GEODIST points Palermo Catania",
				"GEODIST points Palermo Catania km",
			},
			expect: []interface{}{float64(1), float64(1), float64(166274.144), float64(166.2741)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i, cmd := range tc.cmds {
				result, err := exec.FireCommandAndReadResponse(conn, cmd)
				assert.Nil(t, err)
				assert.Equal(t, tc.expect[i], result, "Value mismatch for cmd %s", cmd)
			}
		})
	}
}

func TestGeoHash(t *testing.T) {
	exec := NewWebsocketCommandExecutor()
	conn := exec.ConnectToServer()
	defer conn.Close()

	testCases := []struct {
		name   string
		cmds   []string
		expect []interface{}
	}{
		{
			name:   "GEOHASH with wrong number of arguments",
			cmds:   []string{"GEOHASH points"},
			expect: []interface{}{"ERR wrong number of arguments for 'geohash' command"},
		},
		{
			name: "GEOHASH with non-existent key",
			cmds: []string{
				"GEOHASH nonexistent member1",
			},
			expect: []interface{}{"ERR no such key"},
		},
		{
			name: "GEOHASH with existing key but missing member",
			cmds: []string{
				"GEOADD points -74.0060 40.7128 NewYork",
				"GEOHASH points missingMember",
			},
			expect: []interface{}{float64(1), map[string]interface{}{"missingMember": "nil"}},
		},
		{
			name: "GEOHASH for single member",
			cmds: []string{
				"GEOHASH points NewYork",
			},
			expect: []interface{}{map[string]interface{}{"NewYork": "dr5regw3pp"}},
		},
		{
			name: "GEOHASH for multiple members",
			cmds: []string{
				"GEOADD points -118.2437 34.0522 LosAngeles",
				"GEOHASH points NewYork LosAngeles",
			},
			expect: []interface{}{float64(1), map[string]interface{}{"LosAngeles": "9q5ctr186n", "NewYork": "dr5regw3pp"}},
		},
		{
			name: "GEOHASH with a key of wrong type",
			cmds: []string{
				"SET points somevalue",
				"GEOHASH points member1",
			},
			expect: []interface{}{"OK", "WRONGTYPE Operation against a key holding the wrong kind of value"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i, cmd := range tc.cmds {
				result, err := exec.FireCommandAndReadResponse(conn, cmd)
				assert.Nil(t, err, "Unexpected error for cmd: %s", cmd)
				assert.Equal(t, tc.expect[i], result, "Value mismatch for cmd: %s", cmd)
			}
		})
	}
}
