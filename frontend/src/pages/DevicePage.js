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

const DevicePage = () => {
  const [selectedDevice, setSelectedDevice] = useState(null);
  const [devicesArr, setDevicesArr] = useState([]);
  const [loading, setLoading] = useState(false);
  const toast = useToast();

  useEffect(() => {
    const getDevices = async () => {
      setLoading(true);
      try {
        var response = await axios.get(`${API_BASE_PATH}/devices/`);
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
      setDevicesArr(response.data);
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
        <Text ml="10px">{name}</Text>
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
              value={selectedDevice?.id || devicesArr[0]?.id}
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
        onClick={() => console.log(selectedDevice)}
      >
        OK
      </Button>
    </VStack>
  );
};

export default DevicePage;
