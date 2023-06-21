use std::{env, path::PathBuf};
fn main()->Result<(),Box<dyn std::error::Error>>{

   let out_dir = PathBuf::from(env::var("OUT_DIR").unwrap());
   // println!("{}",format!("Out dir is {}", out_dir.display()));
   tonic_build::configure()
   .build_server(true)
   //.out_dir("./src")
   .file_descriptor_set_path(out_dir.join("user_descriptor.bin"))
   .compile(&["./proto/user.proto"],&["proto"])//.compile(&["../proto/user.proto"],&["../proto"])//first argument is .proto file path, second argument is the folder.
   .unwrap_or_else(|e| panic!("protobuf compile error: {}", e));
   
   Ok(())
}