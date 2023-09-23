import { Flex, Text } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import axios from "axios";
import { API_BASE_PATH } from "./constants";

function App() {
  useEffect(() => {
    const setUp = async () => {
      const response = await axios.get(`${API_BASE_PATH}/`)
      console.log(response)
    }
    setUp();
  }, [])

  return (
    <Flex>
      <Text>
        test
      </Text>
      <Text>
        test
      </Text>
    </Flex>
  );
}

export default App;
