import {
  Center,
  Flex,
  Text,
  Heading,
  Spacer,
  VStack,
  Button,
  HStack,
  useToast,
  Spinner,
} from "@chakra-ui/react";
import { Radio, RadioGroup } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import axios from "axios";
import { API_BASE_PATH } from "../constants";

const EmojiBasedOnType = ({ type }) => {
  let returnVal = "";
  switch (type) {
    case "Computer":
      returnVal = "💻";
      break;
    case "Smartphone":
      returnVal = "📱";
      break;
    case "Speaker":
      returnVal = "🔊";
      break;
  }
  return <Text>{returnVal}</Text>;
};

const DevicePage = ({ setShowDevice }) => {
  const [selectedDevice, setSelectedDevice] = useState(null);
  const [devicesArr, setDevicesArr] = useState([]);
  const [loading, setLoading] = useState(false);
  const toast = useToast();

  useEffect(() => {
    const getDevices = async () => {
      setLoading(true);
      try {
        var response = await axios.get(`${API_BASE_PATH}/devices/`);
        console.log(response);
      } catch (error) {
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
      setDevicesArr(response.data.devices);
      setLoading(false);
    };
    getDevices();
  }, []);

  const ListItem = ({ id, name, type }) => {
    return (
      <Radio
        alignItems="center"
        py="10px"
        colorScheme="green"
        value={id}
        w="100%"
      >
        <HStack spacing="0">
          <EmojiBasedOnType type={type} />
          <Text ml="10px">{name}</Text>
        </HStack>
      </Radio>
    );
  };

  const handleChange = (id) => {
    for (let i = 0; i < devicesArr.length; i++) {
      if (devicesArr[i].id === id) {
        setSelectedDevice(devicesArr[i]);
        break;
      }
    }
  };

  const handleSelectDevice = () => {
    const selectDevice = async () => {
      setLoading(true);
      try {
        // console.log(selectedDevice);
        var response = await axios.post(`${API_BASE_PATH}/devices`, {
          id: selectedDevice.id,
          name: selectedDevice.name,
          type: selectedDevice.type,
        });
      } catch (error) {
        // console.log(error);
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
      if (response.status !== 200) {
        toast({
          title: "Error",
          description: "Something went wrong",
          status: "error",
          duration: 8000,
          isClosable: true,
        });
        setLoading(false);
      }
      if (response.status === 200) {
        setLoading(false);
        setShowDevice(false);
      }
    };
    selectDevice();
  };

  return (
    <VStack w="80vw" maxWidth="450px">
      <Center
        flexDirection="column"
        w="100%"
        bg="gray.600"
        p="20px"
        borderRadius="20px"
      >
        <Heading textAlign="center" px="1em" mb="20px" lineHeight="1">
          Choose your device
        </Heading>
        {loading ? (
          <Spinner colorScheme="green" />
        ) : (
          <Flex
            flexDirection="column"
            alignItems="center"
            justifyContent="center"
            w="100%"
          >
            <RadioGroup
              name="device"
              onChange={handleChange}
              value={selectedDevice?.id}
              maxH="250px"
              overflowY="auto"
            >
              {devicesArr.map((device, index) => {
                return <ListItem {...device} key={index} />;
              })}
            </RadioGroup>
          </Flex>
        )}
      </Center>
      <Button
        colorScheme="green"
        isDisabled={selectedDevice === null}
        w="50%"
        onClick={handleSelectDevice}
      >
        Select
      </Button>
    </VStack>
  );
};

export default DevicePage;
