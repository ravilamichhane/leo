function App(props) {
  console.log("APP rendered", props);
  return (
    <div>
      <h1>タイトル: {props.Name}</h1>
    </div>
  );
}

export default App;
