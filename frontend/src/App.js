import {
  Button,
  Center,
  Flex,
  Text,
  useToast,
  Spinner,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import axios from "axios";
import { API_BASE_PATH } from "./constants";
import DynamicPage from "./pages/DynamicPage";
import DevicePage from "./pages/DevicePage";

const App = () => {
  const toast = useToast();
  const [showLogin, setShowLogin] = useState(false);
  const [showDevice, setShowDevice] = useState(false);

  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const setUp = async () => {
      setLoading(true);
      try {
        var response = await axios.get(`${API_BASE_PATH}/`);
        // console.log(response)
      } catch (error) {
        console.log(error);
        toast({
          title: "Error",
          description: error.response?.data?.detail || error.message,
          status: "error",
          duration: 8000,
          isClosable: true,
        });
        setLoading(false);
        return;
      }

      if (response.data.token === null) {
        setShowLogin(true);
        setLoading(false);
        return;
      }
      if (response.data.device === null) {
        setShowDevice(true);
        setLoading(false);
        return;
      }
      setLoading(false);
    };
    setUp();
  }, []);

  return (
    <Center h="100dvh">
      {loading ? (
        <Spinner color="green" size="xl" />
      ) : (
        <DynamicPage
          showDevice={showDevice}
          setShowLogin={setShowLogin}
          setShowDevice={setShowDevice}
          showLogin={showLogin}
        />
      )}
    </Center>
  );
};

export default App;
