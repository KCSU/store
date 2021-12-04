import { Route, Routes as RRoutes } from "react-router";
import { HelloWorld } from "./views/HelloWorld";
import { Settings } from "./views/Settings";

export const Routes : React.FC = () => {
  return <RRoutes>
    <Route path="/" element={<HelloWorld/>}/>
    <Route path="/settings" element={<Settings/>}/>
  </RRoutes>
};