import { AddIcon, SearchIcon } from "@chakra-ui/icons";
import {
    HStack,
    Heading,
    IconButton,
    Image,
    Input,
    Spacer,
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
    console.log(uri)
    try {
      var response = await axios.post(`${API_BASE_PATH}/current_queue/`, {uri: uri});
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
    <HStack p="20px" bg="gray.500" w="100%" gap="25px">
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
      <Spacer />
      <IconButton
        colorScheme="green"
        icon={loading ? <Spinner /> : <AddIcon />}
        onClick={() => {
          handleSubmitSong();
        }}
      />
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
      <Heading>Search</Heading>
      <HStack w="100%" as="form" onSubmit={handleSubmitQuery}>
        <Input
          focusBorderColor="green.100"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="type here"
        />
        <IconButton colorScheme="green" icon={<SearchIcon />} type="submit" />
      </HStack>
      <VStack
        w="100%"
        maxH="400px"
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
