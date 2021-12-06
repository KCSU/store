import { Flex, FlexProps, forwardRef, useColorModeValue } from "@chakra-ui/react"

export type CardProps = FlexProps;

export const Card = forwardRef<CardProps, 'div'>((props, ref) => {
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
    ref={ref}
    {...props}
  >
  </Flex>
});