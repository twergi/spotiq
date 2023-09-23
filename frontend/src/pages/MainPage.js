import { Button, HStack, Spacer, Text } from "@chakra-ui/react";
import { API_BASE_PATH } from "../constants";
import SearchComponent from "../components/SearchComponent";
import QueueComponent from "../components/QueueComponent";


const MainPage = () => {
    return (

        <HStack w="80vw" maxW="1200px" border="1px solid red" h="90dvh" justifyContent="center">

        {/* <SearchComponent/> */}
        {/* <SearchComponent/> */}
        <QueueComponent/>
        </HStack>
    );
  };

export default MainPage