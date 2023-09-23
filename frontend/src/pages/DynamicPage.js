import LoginPage from "./LoginPage";
import DevicePage from "./DevicePage";

const DynamicPage = ({ showLogin, showDevice }) => {
    if (showLogin === true) {
      return <LoginPage />;
    }
    if (showDevice === true) {
      return <DevicePage />;
    }
  };

export default DynamicPage