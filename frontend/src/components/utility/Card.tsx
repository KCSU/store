import { Flex, FlexProps, useColorModeValue } from "@chakra-ui/react"

export const Card: React.FC<FlexProps> = ({children, ...rest}) => {
  let bg = useColorModeValue('white', 'gray.700');
  return <Flex
    p="22px"
    flexDir="column"
    width="100%"
    pos="relative"
    minW="0px"
    wordBreak="break-word"
    bgClip="border-box"
    bg={bg}
    boxShadow="0px 3.6px 5.5px rgba(0,0,0,0.02)"
    borderRadius="15px"
    {...rest}
  >
    {children}
  </Flex>
}