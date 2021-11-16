// import './App.css'

import { Container, Flex, useBreakpointValue } from "@chakra-ui/react"
import { useState } from "react"
import { Sidebar } from "./components/navigation/Sidebar"
import { HelloWorld } from "./views/HelloWorld"

const smVariant = { navigation: 'drawer', navigationButton: true }
const mdVariant = { navigation: 'sidebar', navigationButton: false }

function App() {
  const [isSidebarOpen, setSidebarOpen] = useState(false);
  const variants = useBreakpointValue({ base: smVariant, md: mdVariant })

  const toggleSidebar = () => setSidebarOpen(!isSidebarOpen)

  return (
    <Flex height="100vh">
      <Sidebar variant={variants?.navigation as 'drawer' | 'sidebar'}
        isOpen={isSidebarOpen}
        onClose={toggleSidebar}/>
      <HelloWorld/>
    </Flex>
  )
}

export default App
