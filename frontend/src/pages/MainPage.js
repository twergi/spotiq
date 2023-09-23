import { VStack } from "@chakra-ui/react";
import { useState } from "react";
import CurrentlyPlaying from "../components/CurrentlyPlaying";
import QueueComponent from "../components/QueueComponent";
import SearchComponent from "../components/SearchComponent";

const MainPage = () => {
  const [currentlyPlaying, setCurrentlyPlaying] = useState(null);

  return (
    <VStack w="80vw" maxW="1200px" overflowY="none">
      {currentlyPlaying !== null && (
        <CurrentlyPlaying currentlyPlaying={currentlyPlaying} />
      )}
      <QueueComponent setCurrentlyPlaying={setCurrentlyPlaying} />
      <SearchComponent />
    </VStack>
  );
};

export default MainPage;
