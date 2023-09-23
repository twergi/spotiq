import { HStack, Image, Text, VStack } from "@chakra-ui/react";


const CurrentlyPlaying = ({currentlyPlaying}) => {
  return (
    <HStack p="20px" bg="green.400" w="100%" gap="25px">
      <Image src={currentlyPlaying.album.images[1].url} alt="image" boxSize="100px" />
      <VStack align="start">
        <Text>{currentlyPlaying.name}</Text>
        <HStack align="start">
          {currentlyPlaying.artists.map((artist, index) => {
            return (
              <Text key={index} mr="5px">
                {artist.name}
              </Text>
            );
          })}
        </HStack>
      </VStack>
    </HStack>
  )
}

export default CurrentlyPlaying