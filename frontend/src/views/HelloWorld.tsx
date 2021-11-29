import { Button } from "@chakra-ui/button";

interface HelloWorldProps {
  onMenu: () => void;
  hasMenuButton: boolean;
};

export const HelloWorld = ({ onMenu, hasMenuButton } : HelloWorldProps) => {
  return <main>
    <div>Hello, World!</div>
    {
      hasMenuButton && <Button onClick={onMenu}>Open Menu</Button>
    }
  </main>;
};
