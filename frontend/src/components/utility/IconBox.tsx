import { Flex } from "@chakra-ui/react";
import React from "react";

export const IconBox: React.FC<any> = ({children, ...rest}) => {
  return (
    <Flex
      alignItems={"center"}
      justifyContent={"center"}
      borderRadius={"12px"}
      {...rest}
    >
      {children}
    </Flex>
  );
}