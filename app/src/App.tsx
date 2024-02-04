import { useEffect, useState } from "react";
import "./App.css";

function App() {
  const [data, setData] = useState(undefined);

  useEffect(() => {
    async function fetchData() {
      const res = await fetch("https://jsonplaceholder.typicode.com/users/1");
      if (res.ok) {
        const json = await res.json();
        setData(json);
      }
    }
    fetchData();
  }, []);

  return <div>{JSON.stringify(data)}</div>;
}

export default App;
