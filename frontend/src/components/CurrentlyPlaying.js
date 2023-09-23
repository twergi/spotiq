import { HStack, Image, Text, VStack } from "@chakra-ui/react";

const CurrentlyPlaying = ({ currentlyPlaying }) => {
  return (
    <HStack p="20px" bg="green.400" w="100%" gap="25px" overflowX="auto">
      <Image
        src={currentlyPlaying?.album?.images[1]?.url}
        alt="image"
        boxSize="100px"
      />
      <VStack align="start">
        <HStack>
          <Text>Name:</Text>
          <Text bg="green.800" p="5px" borderRadius="5px">
            {currentlyPlaying?.name}
          </Text>
        </HStack>
        <HStack align="start">
          <Text>Artists:</Text>
          {currentlyPlaying?.artists?.map((artist, index) => {
            return (
              <Text
                key={index}
                mr="5px"
                bg="green.700"
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

export default CurrentlyPlaying;
