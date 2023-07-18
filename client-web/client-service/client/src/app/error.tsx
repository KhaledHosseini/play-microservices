'use client'
import React from "react";

export default function Error(
    {error}: 
    {error: Error}
    ) {
    return (
      <div className="grid h-screen px-4 bg-white place-content-center">
        <div className="text-center">
          <h1 className="font-black text-gray-200 text-9xl">error</h1>
        </div>
      </div>
    );
  }