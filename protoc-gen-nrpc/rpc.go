package main

import (
	"text/template"

	pb "github.com/LilithGames/nevent/proto"
	npb "github.com/LilithGames/nrpc/proto"

	"github.com/golang/protobuf/proto"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type rpc struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
	tpl *template.Template
}

func NewRpc() pgs.Module {
	return &rpc{ModuleBase: &pgs.ModuleBase{}}
}

func (it *rpc) InitContext(c pgs.BuildContext) {
	it.ModuleBase.InitContext(c)
	it.ctx = pgsgo.InitContext(c.Parameters())
	tpl := template.New("rpc").Funcs(map[string]interface{}{
		"package": it.ctx.PackageName,
		"name":    it.ctx.Name,
		"default": func(value interface{}, defaultValue interface{}) interface{} {
			switch v := value.(type) {
			case string:
				if v == "" {
					return defaultValue
				}
			default:
				panic("default: unknown type")
			}
			return value
		},

		"options": func(node pgs.Node) *pb.EventOption {
			var opt interface{}
			var err error
			switch n := node.(type) {
			case pgs.File:
				opt, err = proto.GetExtension(n.Descriptor().GetOptions(), pb.E_Foptions)
			case pgs.Service:
				opt, err = proto.GetExtension(n.Descriptor().GetOptions(), pb.E_Soptions)
			case pgs.Method:
				opt, err = proto.GetExtension(n.Descriptor().GetOptions(), pb.E_Moptions)
			default:
				panic("node options not supported")
			}
			if err != nil {
				return new(pb.EventOption)
			}
			return opt.(*pb.EventOption)
		},

		"nrpcoptions": func(node pgs.Node) *npb.NRPCOption {
			var opt interface{}
			var err error
			switch n := node.(type) {
			case pgs.Service:
				opt, err = proto.GetExtension(n.Descriptor().GetOptions(), npb.E_Soptions)
			case pgs.Method:
				opt, err = proto.GetExtension(n.Descriptor().GetOptions(), npb.E_Moptions)
			default:
				panic("node options not supported")
			}
			if err != nil {
				return new(npb.NRPCOption)
			}
			return opt.(*npb.NRPCOption)
		},
	})
	it.tpl = template.Must(tpl.Parse(tmpl))
}

func (it *rpc) Name() string {
	return "rpc"
}

func (it *rpc) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		it.generate(f)
	}
	return it.Artifacts()
}

func (it *rpc) generate(f pgs.File) {
	name := it.ctx.OutputPath(f).SetExt(".nrpc.go").String()
	f.Services()
	it.AddGeneratorTemplateFile(name, it.tpl, f)
}
