import { Button, Center, Heading, Image, VStack } from "@chakra-ui/react";
import axios from "axios";
import { useCallback } from "react";
import GoogleLogin, {
  GoogleLoginResponse,
  GoogleLoginResponseOffline,
} from "react-google-login";
import { useLocation, useNavigate } from "react-router-dom";
import { api } from "../config/api";
import { useLogin } from "../hooks/useLogin";
import ravenImg from "../img/raven.png";


export function Login() {
  const location = useLocation();
  const mutation = useLogin();
  const navigate = useNavigate();
  const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;
  const from: string = location.state?.from ?? "/";
  const onLogin = useCallback((response: GoogleLoginResponse | GoogleLoginResponseOffline) => {
    if (!("tokenId" in response)) return;
    mutation.mutateAsync(response.tokenId).then(() => {
      navigate(from);
    })
  }, [mutation])
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
              isLoading={disabled}
              size="lg"
              colorScheme="brand"
              leftIcon={
                <Image
                  src={ravenImg}
                  alt="Raven Logo"
                  width={10}
                  filter="invert()"
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
