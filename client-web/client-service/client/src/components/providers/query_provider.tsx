"use client";

import React, {ReactNode} from "react";
import { QueryClientProvider, QueryClient } from "@tanstack/react-query";
import { useState } from "react";

interface Props {
  children: ReactNode;
}

export const QueryProvider = ({children}: Props) => {
  const queryClientOptions = {
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  };
  const [client] = useState(()=> new QueryClient(queryClientOptions));

  return <QueryClientProvider client={client} contextSharing={true}>{children}</QueryClientProvider>;
}
