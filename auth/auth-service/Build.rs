use std::{env,path::PathBuf};
fn main()->Result<(),Box<dyn std::error::Error>>{
   env::set_var("OUT_DIR", "./src/proto");
   let out_dir = PathBuf::from(env::var("OUT_DIR").unwrap());
   tonic_build::configure()
   .build_server(true)
   .out_dir(out_dir.clone())
   .file_descriptor_set_path(out_dir.join("user_descriptor.bin"))
   .compile(&["./proto/user.proto"],&["proto"])//.compile(&["../proto/user.proto"],&["../proto"])//first argument is .proto file path, second argument is the folder.
   .unwrap_or_else(|e| panic!("protobuf compile error: {}", e));
   
   Ok(())
}