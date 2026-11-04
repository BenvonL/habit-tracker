import React, { useState } from "react";
import { toLocalISODate } from "../../date.js";
import close from "../../assets/icons/remove.svg";
import "./AddActivity.scss";

const AddActivity = ({ habit_id, unit, setAddOpen, getDates }) => {
  const [error, setError] = useState("");

  const incrementCount = (e) => {
    e.preventDefault();
    const count = Number(e.target.count.value);
    if (!Number.isInteger(count) || count < 1) {
      setError("Please enter a whole number of at least 1.");
      return;
    }
    window.backend
      .NewDate(toLocalISODate(), count, habit_id)
      .then((date) =>
        window.backend.MySQLRepository.AddCountFromJSON(JSON.stringify(date))
      )
      .then(() => {
        getDates();
        setAddOpen(false);
      })
      .catch((err) => {
        console.error(err);
        setError(`Could not save: ${err}`);
      });
  };

  return (
    <div className="add-activity">
      <img
        className="close"
        src={close}
        alt="close"
        onClick={() => {
          setAddOpen(false);
        }}
      />
      <h1 className="title">Log Activity</h1>
      <form className="form" onSubmit={incrementCount}>
        <label htmlFor="count">How many {unit} did you complete?</label>
        <input
          className="form-field"
          type="number"
          id="count"
          name="count"
          min="1"
          step="1"
          required
          onChange={() => setError("")}
        />
        {error && <p className="form__error">{error}</p>}
        <div className="form__btn-container">
          <button className="form__btn" type="submit">
            submit
          </button>
          <button
            className="form__btn form__btn--reset"
            onClick={(e) => {
              e.preventDefault();
              setAddOpen(false);
            }}
          >
            cancel
          </button>
        </div>
      </form>
    </div>
  );
};

export default AddActivity;
