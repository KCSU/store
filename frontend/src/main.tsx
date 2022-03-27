import React from "react";
import ReactDOM from "react-dom";
import "@fontsource/montserrat/400.css";
import "@fontsource/montserrat/700.css";
// import './index.css'
import App from "./App";
import { ChakraProvider } from "@chakra-ui/react";
import { BrowserRouter } from "react-router-dom";
import theme from "./config/theme";
import "./config/datetime";
import { QueryCache, QueryClient, QueryClientProvider } from "react-query";
import { ReactQueryDevtools } from 'react-query/devtools'
import axios from "axios";
import { AuthProvider } from "./components/utility/AuthProvider";

// TODO: move to config
const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError(error) {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        queryClient.invalidateQueries('authUser');
      }
    }
  })
})

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        <ChakraProvider theme={theme}>
          <AuthProvider>
            <App />
          </AuthProvider>
        </ChakraProvider>
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById("root")
);
