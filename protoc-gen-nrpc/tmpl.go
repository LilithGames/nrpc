package main

const tmpl = `// This code was autogenerated from nrpc, do not edit.

package {{ package . }}

import (
	"context"
	"fmt"

	"github.com/LilithGames/nevent"
	"github.com/LilithGames/nrpc"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

{{- $foptions := options . }}
{{- $fsubject := default $foptions.Subject (package .) }}
{{- range .Services }}
{{- $svc := .Name }}
{{- $soptions := options . }}
{{- $nrpcoptions := nrpcoptions .}}
{{- $isnrpcsvr := $nrpcoptions.Nrpc }}
{{- $ssubject := default $soptions.Subject $svc }}
{{- $bsubject := printf "%s.%s" $fsubject $ssubject }}
{{- if $isnrpcsvr }}
const (
     {{$svc}}NServiceName = "NRPC4{{$svc}}"
)

type {{$svc}}NInterface interface {
    {{- range .Methods }}
    {{- $nrpcoptions := nrpcoptions .}}
    {{- $isnrpcfun := $nrpcoptions.Nrpc }}
    {{- if $isnrpcfun }}
	{{ name . }}(ctx context.Context, m *{{ name .Input }}) (*{{ name .Output }}, error)
    {{- end }}
    {{- end }}
}

func Register{{$svc}}(s *nrpc.Server, in {{$svc}}NInterface, defaultSubNum int, opts ...nevent.ListenOption) error {
    {{- range .Methods }}
    {{- $oname := name .Output }}
    {{- $moptions := options . }}
    {{- $msubject := default $moptions.Subject (name .) }}
    {{- $subject := printf "nrpc.%s.%s.%s.%s" $svc $fsubject $ssubject $msubject }}
    {{- $nrpcoptions := nrpcoptions .}}
    {{- $isnrpcfun := $nrpcoptions.Nrpc }}

    {{- if $isnrpcfun }}
	GenEh{{ name . }} := func(ctx context.Context, m *nats.Msg) (interface{}, error) {
		data := new({{name .Input }})
		err := proto.Unmarshal(m.Data, data)
		if err != nil {
			return nil, fmt.Errorf("server unmarshal ask: %w", err)
		}

        if err := data.Validate; err != nil {
            return nil, fmt.Errorf("req data validate fail: %w", err)
        }

		resp, err := in.{{ name . }}(ctx, data)
		if err != nil {
			return nil, err
		}

        if err := resp.Validate(); err != nil {
            return nil, fmt.Errorf("rsp data validate fail: %w", err)
        }

		bs, err := proto.Marshal(resp)
		if err != nil {
			return nil, fmt.Errorf("server marshal answer: %w", err)
		}
		return bs, nil
	}

	{{ name .}}Err := s.RegisterEventHandler("{{ $subject }}", defaultSubNum, GenEh{{ name . }}, opts...)
    if {{name .}}Err != nil {
        return {{name .}}Err
    }
    {{- end}}

    {{- end }}
    return nil
}

// generete for nrpc client
type {{ $svc }}NClientImpl struct {
	nc *nevent.Client
    opt []nevent.EmitOption
}

type {{ $svc }}NClient interface {
{{- range .Methods }}
{{- $moptions := options . }}
{{- $msubject := default $moptions.Subject (name .) }}
{{- $subject := printf "nrpc.%s.%s.%s.%s" $svc $fsubject $ssubject $msubject }}
{{- $nrpcoptions := nrpcoptions .}}
{{- $isnrpcfun := $nrpcoptions.Nrpc }}
{{- if $isnrpcfun }}
    {{ name . }}(ctx context.Context, e *{{ name .Input }}, opts ...grpc.CallOption) (*{{ name .Output }}, error)
{{- end }}
{{- end }}
}

func New{{ $svc }}NClient(nc *nevent.Client, opt ...nevent.EmitOption) {{ $svc }}NClient {
    opList := make([]nevent.EmitOption, 0)
    opList = append(opList, opt ...)
	return &{{ $svc }}NClientImpl{nc: nc, opt: opList} 
}

{{- range .Methods }}
{{- $moptions := options . }}
{{- $msubject := default $moptions.Subject (name .) }}
{{- $subject := printf "nrpc.%s.%s.%s" $fsubject $ssubject $msubject }}
{{- $nrpcoptions := nrpcoptions .}}
{{- $isnrpcfun := $nrpcoptions.Nrpc }}
{{- if $isnrpcfun }}

func (nCli *{{ $svc }}NClientImpl){{ name . }}(ctx context.Context, e *{{ name .Input }}, opts ...grpc.CallOption) (*{{ name .Output }}, error) {
	msg := nats.NewMsg("{{ $subject }}")
    if err := e.Validate(); err != nil {
        return nil, fmt.Errorf("ask event invalidate fail %w", vaildErr)
    }
	data, err := proto.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("ask marshal error %w", err)
	}
	msg.Data = data
	resp, err := nCli.nc.Ask(ctx, msg, nCli.opt ...)
	if err != nil {
		return nil, err
	}
	answer := new({{ name .Output }})
	err = proto.Unmarshal(resp, answer)
	if err != nil {
		return nil, fmt.Errorf("answer unmarshal error %w", err)
	}
    if err := resp.Vaildate(); err != nil {
        return nil, fmt.Errorf("rsp validate fail %w", validErr)
    }
	return answer, nil
}
{{- end }}
{{- end }}
{{- end }}

{{- end }}
`
