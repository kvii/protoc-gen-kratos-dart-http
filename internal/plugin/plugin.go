package plugin

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func Generate(request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	registry, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{
		File: request.ProtoFile,
	})
	if err != nil {
		return nil, fmt.Errorf("create proto registry: %w", err)
	}

	var res pluginpb.CodeGeneratorResponse
	for _, p := range request.FileToGenerate {
		fd, err := registry.FindFileByPath(p)
		if err != nil {
			return nil, fmt.Errorf("can not find file by path: %w", err)
		}

		if fd.Messages().Len() > 0 {
			f := new(File)
			message{path: p, fd: fd, f: f}.Generate()

			res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(changeExtension(p, ".dart")),
				Content: proto.String(string(f.Content())),
			})
		}

		if fd.Enums().Len() > 0 {
			f := new(File)
			enum{fd: fd, f: f}.Generate()

			res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(changeExtension(p, ".enum.dart")),
				Content: proto.String(string(f.Content())),
			})
		}

		cg := client{path: p, fd: fd, f: new(File)}
		if cg.Need() {
			cg.Generate()

			res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(changeExtension(p, ".client.dart")),
				Content: proto.String(string(cg.f.Content())),
			})
		}
	}

	res.SupportedFeatures = proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL))
	return &res, nil
}
