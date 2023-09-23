import {
  Center,
  Spinner,
  useToast
} from "@chakra-ui/react";
import axios from "axios";
import { useEffect, useState } from "react";
import { API_BASE_PATH } from "./constants";
import DynamicPage from "./pages/DynamicPage";

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
    <Center mt="20px">
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
