import { Button, Icon } from "@chakra-ui/react";
import { FaArrowLeft } from "react-icons/fa";
import { Link, To } from "react-router-dom";

export interface BackButtonProps {
  to?: To;
}

export const BackButton: React.FC<BackButtonProps> = ({children, to = "/"}) => {
  return (
    <Button
      as={Link}
      // size="sm"
      to={to}
      variant="ghost"
      mb={4}
      leftIcon={<Icon as={FaArrowLeft} />}
    >
      {children}
    </Button>
  );
}
