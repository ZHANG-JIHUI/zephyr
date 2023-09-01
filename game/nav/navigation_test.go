package nav

import (
	"fmt"
	"github.com/arl/gogeo/f32/d3"
	"testing"
)

func TestFindStraightPath(t *testing.T) {
	nav := NewNavigation()
	err := nav.LoadNavMesh("all_tiles_navmesh.bin")
	if err != nil {
		return
	}
	org := d3.Vec3{24.4, 0.0, 9.0}
	dst := d3.Vec3{-24.0, 2.0, 18.0}
	path, err := nav.FindPathStraightPath(org, dst)
	if err != nil {
		return
	}
	fmt.Println(path)
}

func TestCanReach(t *testing.T) {
	nav := NewNavigation()
	_ = nav.LoadNavMesh("all_tiles_navmesh.bin")
	dst := d3.Vec3{-24.0, 2.0, 18.0}
	fmt.Println(nav.CanReach(dst))
}

func BenchmarkNavigation_FindPathStraightPath(b *testing.B) {
	nav := NewNavigation()
	err := nav.LoadNavMesh("all_tiles_navmesh.bin")
	if err != nil {
		return
	}
	org := d3.Vec3{24.4, 0.0, 9.0}
	dst := d3.Vec3{-24.0, 2.0, 18.0}
	for i := 0; i < b.N; i++ {
		_, _ = nav.FindPathStraightPath(org, dst)
	}
}
