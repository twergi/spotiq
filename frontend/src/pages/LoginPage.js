import { Button } from "@chakra-ui/react";
import { API_BASE_PATH } from "../constants";

const LoginPage = () => {
    return (
      <Button
        size="lg"
        colorScheme="green"
        my="auto"
        as="a"
        href={`${API_BASE_PATH}/login`}
      >
        Login
      </Button>
    );
  };

export default LoginPage