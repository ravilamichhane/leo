import React from "react";

const Home = (props) => {
  React.useEffect(() => {
    console.log("|HAHAHA", props);
  });
  return <div>Home {props.name}</div>;
};

export default Home;