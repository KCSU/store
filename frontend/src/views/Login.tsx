import { Button, Center, Heading, Image, useColorModeValue, VStack } from "@chakra-ui/react";
import axios from "axios";
import { useCallback } from "react";
import GoogleLogin, {
  GoogleLoginResponse,
  GoogleLoginResponseOffline,
} from "react-google-login";
import { Navigate, useLocation, useNavigate } from "react-router-dom";
import { api } from "../config/api";
import { useAuthUser } from "../hooks/useAuthUser";
import { useLogin } from "../hooks/useLogin";
import ravenImg from "../img/raven.png";


export function Login() {
  const location = useLocation();
  const mutation = useLogin();
  const {data: user, isError} = useAuthUser();
  const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;
  const from: string = location.state?.from.pathname ?? "/";
  const onLogin = useCallback((response: GoogleLoginResponse | GoogleLoginResponseOffline) => {
    if (!("tokenId" in response)) return;
    mutation.mutate(response.tokenId);
  }, [mutation]);
  const filter = useColorModeValue('invert()', undefined);
  if (user && !isError) {
    return <Navigate to={from} />
  }
  return (
    <VStack mt={40}>
      <Heading as="h1" size="xl">
        KiFoMaSy
      </Heading>
      <GoogleLogin
        clientId={clientId}
        onSuccess={onLogin}
        cookiePolicy="single_host_origin"
        hostedDomain="cam.ac.uk"
        render={({ onClick, disabled }) => {
          return (
            <Button
              onClick={onClick}
              isLoading={disabled || mutation.isLoading}
              size="lg"
              colorScheme="brand"
              leftIcon={
                <Image
                  src={ravenImg}
                  alt="Raven Logo"
                  width={10}
                  filter={filter}
                />
              }
            >
              Log in with Raven
            </Button>
          );
        }}
      />
    </VStack>
  );
}
