import {
  Button,
  Heading,
  Image,
  useColorModeValue,
  VStack,
} from "@chakra-ui/react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuthUser } from "../hooks/useAuthUser";
import ravenImg from "../img/raven.png";

export function LoginView() {
  const location = useLocation();
  const { data: user, isError, isLoading } = useAuthUser();
  const from: string = (location.state as any)?.from ?? "/";
  const filter = useColorModeValue("invert()", undefined);
  const redirectUrl = import.meta.env.VITE_API_BASE_URL + 'oauth/redirect';
  if (user && !isError && !isLoading) {
    return <Navigate to={from} />;
  }
  return (
    <VStack mt={40}>
      <Heading as="h1" size="xl">
        KiFoMaSy
      </Heading>

      <Button
        as="a" href={redirectUrl}
        size="lg"
        colorScheme="brand"
        leftIcon={
          <Image src={ravenImg} alt="Raven Logo" width={10} filter={filter} />
        }
      >
        Log in with Raven
      </Button>
    </VStack>
  );
}
