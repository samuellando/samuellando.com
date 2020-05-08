import React, { useState } from "react";
import { Auth } from "aws-amplify";
import "./Login.css";

export default function Login(props) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  function validateForm() {
    return email.length > 0 && password.length > 0;
  }

  async function handleSubmit(event) {
    event.preventDefault();
  
    try {
      await Auth.signIn(email, password);
        alert("Logged in");
      } catch (e) {
        alert(e.message);
      }
  }

  return (
    <div className="Login">
      <form onSubmit={handleSubmit}>
          <label>Email</label>
          <input
            type="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
          />
          <label>Password</label>
          <input
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
          />
        <input type="submit"/>
      </form>
    </div>
  );
}

