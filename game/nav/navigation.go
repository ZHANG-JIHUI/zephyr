package nav

import (
	"github.com/arl/go-detour/detour"
	"github.com/arl/gogeo/f32/d3"
	"github.com/pkg/errors"
	"os"
)

const (
	MaxPolys = 256
	MaxNodes = 1024
)

type Navigation struct {
	path           string
	navmesh        *detour.NavMesh
	query          *detour.NavMeshQuery
	defaultExtents d3.Vec3                     // default extents vector for the nearest polygon query
	defaultFilter  *detour.StandardQueryFilter // default query filter
}

func NewNavigation() *Navigation {
	nav := &Navigation{
		defaultExtents: d3.NewVec3XYZ(2, 4, 2),
		defaultFilter:  detour.NewStandardQueryFilter(),
	}
	return nav
}

func (slf *Navigation) LoadNavMesh(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	slf.navmesh, err = detour.Decode(fp)
	if err != nil {
		return err
	}
	st, query := detour.NewNavMeshQuery(slf.navmesh, MaxNodes)
	if detour.StatusFailed(st) {
		return errors.Errorf("query creation failed with status 0x%x", st)
	}
	slf.query = query
	return nil
}

// FindPathStraightPath 查询路径
func (slf *Navigation) FindPathStraightPath(org, dst d3.Vec3) ([]d3.Vec3, error) {
	// query
	query := slf.query
	// get org polygon reference
	st, orgRef, org := query.FindNearestPoly(org, slf.defaultExtents, slf.defaultFilter)
	if detour.StatusFailed(st) {
		return nil, errors.Errorf("couldn't find nearest poly of %v, status: 0x%x", org, st)
	}
	if !slf.navmesh.IsValidPolyRef(orgRef) {
		return nil, errors.Errorf("orgRef %d is not a valid poly ref", orgRef)
	}

	// get dst polygon reference
	st, dstRef, dst := query.FindNearestPoly(dst, slf.defaultExtents, slf.defaultFilter)
	if detour.StatusFailed(st) {
		return nil, errors.Errorf("couldn't find nearest poly of %v, status: 0x%x", dst, st)
	}
	if !slf.navmesh.IsValidPolyRef(dstRef) {
		return nil, errors.Errorf("dstRef %d is not a valid poly ref", dstRef)
	}

	path := make([]detour.PolyRef, MaxPolys)
	pathCount, st := query.FindPath(orgRef, dstRef, org, dst, slf.defaultFilter, path)
	if detour.StatusFailed(st) {
		return nil, errors.Errorf("query.FindPath failed with 0x%x", st)
	}
	straightPath := make([]d3.Vec3, MaxPolys)
	for i := range straightPath {
		straightPath[i] = d3.NewVec3()
	}
	straightPathFlags := make([]uint8, MaxPolys)
	straightPathRefs := make([]detour.PolyRef, MaxPolys)
	straightPathCount, st := query.FindStraightPath(org, dst, path[:pathCount], straightPath, straightPathFlags, straightPathRefs, 0)
	if detour.StatusFailed(st) {
		return nil, errors.Errorf("query.FindStraightPath failed with 0x%x", st)
	}
	if (straightPathFlags[0] & detour.StraightPathStart) == 0 {
		return nil, errors.Errorf("straightPath start is not flagged StraightPathStart")
	}
	if (straightPathFlags[straightPathCount-1] & detour.StraightPathEnd) == 0 {
		return nil, errors.Errorf("straightPath end is not flagged StraightPathEnd")
	}
	return straightPath[:straightPathCount], nil
}

// CanReach 是否可达
func (slf *Navigation) CanReach(dst d3.Vec3) bool {
	st, dstRef, _ := slf.query.FindNearestPoly(dst, slf.defaultExtents, slf.defaultFilter)
	if detour.StatusFailed(st) || !slf.navmesh.IsValidPolyRef(dstRef) {
		return false
	}
	return true
}
