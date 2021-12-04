import { IconButton } from "@chakra-ui/button";
import { useColorMode } from "@chakra-ui/react";
import { FaHamburger, FaMoon, FaSun } from "react-icons/fa";

interface HelloWorldProps {
  onMenu: () => void;
  hasMenuButton: boolean;
};

export const HelloWorld = ({ onMenu, hasMenuButton } : HelloWorldProps) => {
  const { colorMode, toggleColorMode } = useColorMode();

  return <main>
    <div>Hello, World!</div>
    <IconButton onClick={toggleColorMode} aria-label="toggle dark">
      {colorMode == 'light' ? <FaMoon/> : <FaSun/>}
    </IconButton>
    {
      hasMenuButton && <IconButton aria-label="open menu" onClick={onMenu}><FaHamburger/></IconButton>
    }
  </main>;
};
