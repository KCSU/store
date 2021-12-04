// import './App.css'

import { Flex, useBreakpointValue, useColorModeValue } from "@chakra-ui/react"
import { useState } from "react"
import { Sidebar } from "./components/navigation/Sidebar"
import { HelloWorld } from "./views/HelloWorld"

const smVariant = { navigation: 'drawer', navigationButton: true }
const mdVariant = { navigation: 'sidebar', navigationButton: false }

function App() {
  const [isSidebarOpen, setSidebarOpen] = useState(false);
  const variants = useBreakpointValue({ base: smVariant, md: mdVariant })

  const toggleSidebar = () => setSidebarOpen(!isSidebarOpen)
  const bg = useColorModeValue("gray.50", "gray.800");

  return (
    <Flex height="100vh" bg={bg}>
      <Sidebar variant={variants?.navigation as 'drawer' | 'sidebar'}
        isOpen={isSidebarOpen}
        onClose={toggleSidebar}/>
      {/* TODO: Routes */}
      <HelloWorld onMenu={toggleSidebar}
        hasMenuButton={!!variants?.navigationButton}/>
    </Flex>
  )
}

export default App
