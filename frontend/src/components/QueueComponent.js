import {
    HStack,
    Text,
    Image,
    Input,
    VStack,
    useToast,
    Spinner,Spacer, Center, Heading,
  } from "@chakra-ui/react";
  import { useState } from "react";
  import { IconButton } from "@chakra-ui/react";
  import { SearchIcon, AddIcon} from "@chakra-ui/icons";
  import axios from "axios";
  import { API_BASE_PATH } from "../constants";
  
  const ResultsItem = ({ uri, name, artistsNamesArr, imageUrl }) => {
    return (
      <HStack p="20px" bg="gray.400" w="100%" gap="25px">
        <Image src={imageUrl.url} alt="image" boxSize="100px" />
        <VStack align="start">
          <Text>{name}</Text>
          <HStack align="start">
            {artistsNamesArr.map((artist, index) => {
              return (
                <Text key={index} mr="5px">
                  {artist.name}
                </Text>
              );
            })}
          </HStack>
        </VStack>
        <Spacer/>
        <IconButton colorScheme="green" icon={<AddIcon/>}/>
      </HStack>
    );
  };
  
  const QueueComponent = ({}) => {
    const [query, setQuery] = useState("");
    const [loading, setLoading] = useState(false);
    const [queueResults, setQueueResults] = useState([]);
  
    const toast = useToast();
  
    
    return (
      <VStack w="45%" h="100%">
      <Heading>
          Queue
      </Heading>
        <VStack  w="100%" h="100%"  overflowY={loading ? "none" : "auto"} align="start">
          {loading ? (
            <Spinner mx="auto" color="green" size="lg" />
          ) : (
            queueResults.map((item) => {
              return (
                <ResultsItem
                  key={item.uri}
                  name={item.name}
                  artistsNamesArr={item.artists}
                  imageUrl={item.album.images[0]}
                />
              );
            })
          )}
        </VStack>
      </VStack>
    );
  };
  
  export default QueueComponent;
  