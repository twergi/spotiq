import { AddIcon, CloseIcon, SearchIcon } from "@chakra-ui/icons";
import {
  HStack,
  Heading,
  IconButton,
  Image,
  Input,
  InputGroup,
  InputRightElement,
  Spinner,
  Text,
  VStack,
  useToast
} from "@chakra-ui/react";
import axios from "axios";
import { useState } from "react";
import { API_BASE_PATH } from "../constants";

const ResultsItem = ({ uri, name, artistsNamesArr, imageUrl }) => {
  const toast = useToast();
  const [loading, setLoading] = useState(false);

  const handleSubmitSong = async () => {
    setLoading(true);
    console.log(uri);
    try {
      var response = await axios.post(`${API_BASE_PATH}/current_queue/`, {
        uri: uri,
      });
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
    if (response.status === 200) {
      toast({
        title: "Success",
        status: "success",
        position: "top",
        duration: 8000,
        isClosable: true,
      });
      setLoading(false);
    }
  };

  return (
    <HStack p="20px" bg="gray.700" w="100%" gap="25px" overflowX="auto">
      <IconButton
        bg="green.400"
        icon={loading ? <Spinner /> : <AddIcon />}
        onClick={() => {
          handleSubmitSong();
        }}
      />
      <Image src={imageUrl.url} alt="image" boxSize="100px" />
      <VStack align="start">
      <HStack>
          <Text>Name:</Text>
          <Text bg="gray.900" p="5px" borderRadius="5px">
            {name}
          </Text>
        </HStack>
        <HStack align="center">
          <Text>
            Artists:
          </Text>
          {artistsNamesArr.map((artist, index) => {
            return (
              <Text
                key={index}
                mr="5px"
                bg="gray.800"
                p="5px"
                borderRadius="5px"
              >
                {artist.name}
              </Text>
            );
          })}
        </HStack>
      </VStack>
    </HStack>
  );
};

const SearchComponent = () => {
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [searchResults, setSearchResults] = useState([]);

  const toast = useToast();

  const handleSubmitQuery = (e) => {
    e.preventDefault();
    const lookup = async () => {
      setLoading(true);
      try {
        var response = await axios.get(`${API_BASE_PATH}/search?q=${query}`);
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
      console.log(response.data.items);
      setSearchResults(response.data.items);
      setLoading(false);
    };
    lookup();
  };

  return (
    <VStack w="100%">
      <Heading fontSize="25px" alignSelf="start" ml="18px">Search</Heading>
      <HStack w="100%" as="form" onSubmit={handleSubmitQuery}>
      <InputGroup>
        <Input
          focusBorderColor="green.100"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search for track..."
          />
          {query !== "" &&
          <InputRightElement _hover={{cursor: "pointer"}} onClick={() => {
            setQuery("");
            setSearchResults([]);
          }} children={<CloseIcon/>}/> }
          </InputGroup>
        <IconButton colorScheme="green" icon={<SearchIcon />} type="submit" />

      </HStack>
      <VStack
        w="100%"
        // maxH="400px"
        overflowY={loading ? "none" : "auto"}
        align="start"
      >
        {loading ? (
          <Spinner mx="auto" color="green" size="lg" />
        ) : (
          searchResults.map((item) => {
            return (
              <ResultsItem
                uri={item.uri}
                key={item.uri}
                name={item.name}
                artistsNamesArr={item.artists}
                imageUrl={item.album.images[1]}
              />
            );
          })
        )}
      </VStack>
    </VStack>
  );
};

export default SearchComponent;
