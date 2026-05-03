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
    <HStack bg="green.800" w="100%" gap="25px" overflowX="auto">
      <Image 
        src={imageUrl?.url || ''} 
        alt="Track image" 
        boxSize="100px" 
        fallbackSrc="https://placehold.co/100"
      />
      <VStack align="start">
        <HStack align="center">
          <Text bg="green.600" p="5px" borderRadius="5px">
            {name || 'Unknown track'}
          </Text>
        </HStack>
        <HStack align="center">
          <Text mr="5px" bg="green.700" p="5px" borderRadius="5px">
            {artistsNamesArr?.map(artist => artist?.name).filter(Boolean).join(', ') || 'Unknown artist'}
          </Text>
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
    const response = await axios.get(`${API_BASE_PATH}/current_queue`);

    setQueueResults(response?.data?.queue || []);
    setCurrentlyPlaying(response?.data?.currently_playing || null);
    
    console.log(response);
  } catch (error) {
    console.error("Queue update error:", error);
    
    const errorMessage = error?.response?.data?.detail?.message 
      || error?.message 
      || "Unknown error occurred";
    
    toast({
      title: "Error",
      description: errorMessage,
      status: "error",
      duration: 8000,
      isClosable: true,
    });
    
    setQueueResults([]);
    setCurrentlyPlaying(null);
  } 
};

const QueueComponent = ({ setCurrentlyPlaying }) => {
  const [loading, setLoading] = useState(false);
  const [queueResults, setQueueResults] = useState([]);
  const toast = useToast();

  useEffect(() => {
    let isMounted = true;
    let interval;

    const fetchData = async () => {
      if (!isMounted) return;
      await updateQueue(setQueueResults, toast, setLoading, setCurrentlyPlaying);
    };

    fetchData();

    interval = setInterval(() => {
      if (isMounted) {
        updateQueue(setQueueResults, toast, setLoading, setCurrentlyPlaying);
      }
    }, 5000);

    return () => {
      isMounted = false;
      if (interval) {
        clearInterval(interval);
      }
    };
  }, [setCurrentlyPlaying, toast]);

  const hasQueueResults = queueResults && Array.isArray(queueResults) && queueResults.length > 0;

  return (
    <Accordion w="100%" allowToggle>
      <AccordionItem w="100%" border="none">
        <AccordionButton w="100%" pl="0px">
          <HStack spacing="20px">
            <Heading fontSize="25px" ml="2px">Очередь</Heading>
            <AccordionIcon style={{ scale: "2", marginTop: "7px" }} />
          </HStack>
        </AccordionButton>
        <AccordionPanel p="0" w="100%">
          <VStack
            w="100%"
            overflowY={loading ? "none" : "auto"}
            align="start"
          >
            {loading ? (
              <Spinner mx="auto" color="green" size="lg" />
            ) : !hasQueueResults ? (
              <Text p="20px" color="gray.500" textAlign="center" w="100%">
                Queue is empty
              </Text>
            ) : (
              queueResults.map((item, index) => (
                <ResultsItem
                  key={item?.id || index}
                  name={item?.name}
                  artistsNamesArr={item?.artists}
                  imageUrl={item?.album?.images?.[1]}
                />
              ))
            )}
          </VStack>
        </AccordionPanel>
      </AccordionItem>
    </Accordion>
  );
};

export default QueueComponent;