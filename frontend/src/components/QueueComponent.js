import {
  HStack,
  Heading,
  Image,
  Spinner,
  Text,
  VStack,
  useToast,
} from "@chakra-ui/react";
import axios from "axios";
import { useEffect, useState } from "react";
import { API_BASE_PATH } from "../constants";

import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
} from "@chakra-ui/react";

const ResultsItem = ({ name, artistsNamesArr, imageUrl }) => {
  return (
    <HStack p="20px" bg="green.800" w="100%" gap="25px" overflowX="auto">
      <Image src={imageUrl?.url} alt="image" boxSize="100px" />
      <VStack align="start">
        <HStack align="center">
          <Text>Name:</Text>
          <Text bg="green.600" p="5px" borderRadius="5px">
            {name}
          </Text>
        </HStack>
        <HStack align="center">
          <Text>Artists:</Text>
          {artistsNamesArr?.map((artist, index) => {
            return (
              <Text
                key={index}
                mr="5px"
                bg="green.700"
                p="5px"
                borderRadius="5px"
              >
                {artist?.name}
              </Text>
            );
          })}
        </HStack>
      </VStack>
    </HStack>
  );
};

const updateQueue = async (
  setQueueResults,
  toast,
  setLoading,
  setCurrentlyPlaying
) => {
  try {
    var response = await axios.get(`${API_BASE_PATH}/current_queue/`);
  } catch (error) {
    console.log(error);
    toast({
      title: "Error",
      description: error.response?.data?.detail.message || error.message,
      status: "error",
      duration: 8000,
      isClosable: true,
    });
    setLoading(false);
    return;
  }
  setQueueResults(response.data.queue);
  setCurrentlyPlaying(response.data.currently_playing);
  console.log(response);
};

const QueueComponent = ({ setCurrentlyPlaying }) => {
  const [loading, setLoading] = useState(false);
  const [queueResults, setQueueResults] = useState([]);
  const toast = useToast();

  useEffect(() => {
    updateQueue(setQueueResults, toast, setLoading, setCurrentlyPlaying);
    let interval = setInterval(() => {
      updateQueue(setQueueResults, toast, setLoading, setCurrentlyPlaying);
    }, 5000);
    return () => {
      clearInterval(interval);
    };
  }, []);

  return (
    <Accordion w="100%" allowToggle>
      <AccordionItem w="100%" border="none">
        <AccordionButton w="100%">
          <HStack spacing="20px">
            <Heading fontSize="25px">Queue</Heading>{" "}
            <AccordionIcon style={{ scale: "2", marginTop: "7px" }} />
          </HStack>
        </AccordionButton>
        <AccordionPanel p="0" w="100%">
          <VStack
            w="100%"
            // maxH="400px"
            overflowY={loading ? "none" : "auto"}
            align="start"
          >
            {loading ? (
              <Spinner mx="auto" color="green" size="lg" />
            ) : (
              queueResults?.map((item, index) => {
                return (
                  <ResultsItem
                    key={index}
                    name={item?.name}
                    artistsNamesArr={item?.artists}
                    imageUrl={item?.album?.images[1]}
                  />
                );
              })
            )}
          </VStack>
        </AccordionPanel>
      </AccordionItem>
    </Accordion>
  );
};

export default QueueComponent;
