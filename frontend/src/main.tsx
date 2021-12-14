import React from "react";
import ReactDOM from "react-dom";
import '@fontsource/montserrat/400.css'
import '@fontsource/montserrat/700.css'
// import './index.css'
import App from "./App";
import { ChakraProvider } from "@chakra-ui/react";
import { BrowserRouter } from "react-router-dom";
import theme from "./config/theme";
import "./config/datetime";

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <ChakraProvider theme={theme}>
        <App />
      </ChakraProvider>
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById("root")
);
