import { Button, Spacer, VStack } from "@chakra-ui/react";
import { useState } from "react";
import CurrentlyPlaying from "../components/CurrentlyPlaying";
import QueueComponent from "../components/QueueComponent";
import SearchComponent from "../components/SearchComponent";
import { API_BASE_PATH } from "../constants";

const MainPage = () => {
  const [currentlyPlaying, setCurrentlyPlaying] = useState(null);

  return (
    <>
      <VStack h="95dvh">
        <VStack w="80vw" maxW="1200px" pb="20px" spacing="20px">
          {currentlyPlaying !== null && (
            <CurrentlyPlaying currentlyPlaying={currentlyPlaying} />
          )}
          <QueueComponent setCurrentlyPlaying={setCurrentlyPlaying} />
          <SearchComponent />
          <Spacer />
        </VStack>
        <Spacer/>
        <Button
          as="a"
          href={`${API_BASE_PATH}/logout`}
          alignSelf="flex-end"
          colorScheme="red"
          minW="80px"
          minH="40px"
        >
          Logout
        </Button>
        <Spacer minH="20px" maxH="20px"/>
      </VStack>
    </>
  );
};

export default MainPage;
