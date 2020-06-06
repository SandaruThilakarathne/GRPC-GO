protoc -I greet/ greet/greetpb/greet.proto --go_out=plugins=grpc:greet
protoc -I blog/ blog/blogpb/blog.proto --go_out=plugins=grpc:blog