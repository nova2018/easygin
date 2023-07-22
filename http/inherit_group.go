package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nova2018/goutils"
	"net/http"
)

type Group struct {
	groups []*GroupWithTable
}

type routerConfig struct {
	group        *gin.RouterGroup
	httpMethod   string
	relativePath string
	handlers     []gin.HandlerFunc
}

type GroupWithTable struct {
	routerGroup *gin.RouterGroup
	routerTable *map[string]map[string]*routerConfig
}

func (g *GroupWithTable) Sync() {
	for _, temp := range *g.routerTable {
		for _, routerCfg := range temp {
			if routerCfg.group == g.routerGroup {
				g.routerGroup.Handle(routerCfg.httpMethod, routerCfg.relativePath, routerCfg.handlers...)
			}
		}
	}
}

func (g *Group) Group(relativePath string, handlers ...gin.HandlerFunc) *Group {
	list := make([]*GroupWithTable, 0, len(g.groups))
	for _, group := range g.groups {
		list = append(list, &GroupWithTable{
			routerGroup: group.routerGroup.Group(relativePath, handlers...),
			routerTable: group.routerTable,
		})
	}
	return &Group{groups: list}
}

func (g *Group) each(fn func(*GroupWithTable)) gin.IRoutes {
	for _, gp := range g.groups {
		fn(gp)
	}
	return g
}

func (g *Group) handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return g.each(func(group *GroupWithTable) {
		absPath := goutils.JoinPaths(group.routerGroup.BasePath(), relativePath)
		if _, ok := (*group.routerTable)[httpMethod]; !ok {
			(*group.routerTable)[httpMethod] = make(map[string]*routerConfig, 0)
		}
		(*group.routerTable)[httpMethod][absPath] = &routerConfig{
			group:        group.routerGroup,
			httpMethod:   httpMethod,
			relativePath: relativePath,
			handlers:     handlers,
		}
	})
}

func (g *Group) Sync() {
	g.groups[0].Sync()
}

func (g *Group) Use(handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.each(func(group *GroupWithTable) {
		group.routerGroup.Use(handlerFunc...)
	})
}

func (g *Group) Handle(s string, s2 string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.handle(s, s2, handlerFunc...)
}

func (g *Group) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	g.handle(http.MethodGet, relativePath, handlers...)
	g.handle(http.MethodPost, relativePath, handlers...)
	g.handle(http.MethodPut, relativePath, handlers...)
	g.handle(http.MethodPatch, relativePath, handlers...)
	g.handle(http.MethodHead, relativePath, handlers...)
	g.handle(http.MethodOptions, relativePath, handlers...)
	g.handle(http.MethodDelete, relativePath, handlers...)
	g.handle(http.MethodConnect, relativePath, handlers...)
	g.handle(http.MethodTrace, relativePath, handlers...)
	return g
}

func (g *Group) GET(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.handle(http.MethodGet, s, handlerFunc...)
}

func (g *Group) POST(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.handle(http.MethodPost, s, handlerFunc...)
}

func (g *Group) DELETE(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.handle(http.MethodDelete, s, handlerFunc...)

}

func (g *Group) PATCH(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.handle(http.MethodPatch, s, handlerFunc...)

}

func (g *Group) PUT(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.handle(http.MethodPut, s, handlerFunc...)

}

func (g *Group) OPTIONS(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.handle(http.MethodOptions, s, handlerFunc...)

}

func (g *Group) HEAD(s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.handle(http.MethodHead, s, handlerFunc...)

}

func (g *Group) StaticFile(s string, s2 string) gin.IRoutes {
	return g.each(func(group *GroupWithTable) {
		group.routerGroup.StaticFile(s, s2)
	})
}

func (g *Group) Static(s string, s2 string) gin.IRoutes {
	return g.each(func(group *GroupWithTable) {
		group.routerGroup.Static(s, s2)
	})
}

func (g *Group) StaticFS(s string, system http.FileSystem) gin.IRoutes {
	return g.each(func(group *GroupWithTable) {
		group.routerGroup.StaticFS(s, system)
	})
}

func (g *Group) Match(strings []string, s string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
	return g.each(func(group *GroupWithTable) {
		group.routerGroup.Match(strings, s, handlerFunc...)
	})
}

func (g *Group) StaticFileFS(s string, s2 string, system http.FileSystem) gin.IRoutes {
	return g.each(func(group *GroupWithTable) {
		group.routerGroup.StaticFileFS(s, s2, system)
	})
}

func NewInheritGroup(list ...*gin.RouterGroup) []*Group {
	listNewGroup := make([]*GroupWithTable, 0, len(list))
	routerTable := make(map[string]map[string]*routerConfig, 0)
	for _, g := range list {
		listNewGroup = append(listNewGroup, &GroupWithTable{
			routerGroup: g,
			routerTable: &routerTable,
		})
	}
	listGroup := make([]*Group, 0, len(list))
	for i := range listNewGroup {
		listGroup = append(listGroup, &Group{
			groups: listNewGroup[i:],
		})
	}
	return listGroup
}
