import LoginPage from "./LoginPage";
import DevicePage from "./DevicePage";
import MainPage from "./MainPage"

const DynamicPage = ({ showLogin, showDevice, setShowDevice }) => {
    if (showLogin === true) {
      return <LoginPage />;
    }
    if (showDevice === true) {
      return <DevicePage setShowDevice={setShowDevice}/>;
    }
    return <MainPage />
  };

export default DynamicPage