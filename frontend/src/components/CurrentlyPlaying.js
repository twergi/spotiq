import { HStack, Image, Text, VStack } from "@chakra-ui/react";

const CurrentlyPlaying = ({ currentlyPlaying }) => {
  if (!currentlyPlaying) {
    return (
      <HStack bg="green.400" w="100%" gap="25px" overflowX="auto">
        <Text>No track currently playing</Text>
      </HStack>
    );
  }

  return (
    <HStack bg="green.400" w="100%" gap="25px" overflowX="auto">
      <Image
        src={currentlyPlaying?.album?.images?.[0]?.url || ''}
        alt="image"
        boxSize="100px"
        fallbackSrc="https://placehold.co/100"
      />
      <VStack align="start">
        <HStack>
          <Text bg="green.800" p="5px" borderRadius="5px">
            {currentlyPlaying?.name || 'Unknown track'}
          </Text>
        </HStack>
        <HStack align="start">
          <Text mr="5px" bg="green.700" p="5px"borderRadius="5px">
            {currentlyPlaying?.artists?.map(artist => artist?.name).filter(Boolean).join(', ') || 'Unknown artist'}
          </Text>
        </HStack>
      </VStack>
    </HStack>
  );
};

export default CurrentlyPlaying;