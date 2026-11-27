import React, { useState, useEffect } from "react";
import CalendarHeatmap from "react-calendar-heatmap";
import ReactTooltip from "react-tooltip";
import edit from "../../assets/icons/btn--edit.png";
import add from "../../assets/icons/btn--add.png";
import AddActivity from "../../components/AddActivity";
import EditHabit from "../../components/EditHabit";
import Pomodoro from "../../components/Pomodoro";
import { toLocalISODate, shiftDate, countToLevel } from "../../date.js";
import "./Habit.scss";

const DAYS_SHOWN = 364;

const Habit = ({ habitList, updateHabits }) => {
  const url = window.location.pathname.split("/").pop();
  const [addOpen, setAddOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [habit, setHabit] = useState({});
  const [dates, setDates] = useState([]);

  const today = toLocalISODate();
  const maxCount = Math.max(0, ...dates.map((d) => d.count));

  const getDates = () => {
    window.backend.MySQLRepository.GetAllDates(habit.id)
      .then((response) => {
        // one value per day: start with 0 for every shown day, then
        // overwrite with what is stored in the database
        const countByDate = new Map(
          getRange(DAYS_SHOWN).map((index) => [shiftDate(today, -index), 0])
        );
        (response || []).forEach((entry) => {
          countByDate.set(entry.date, entry.count);
        });
        setDates(
          Array.from(countByDate, ([date, count]) => ({ date, count }))
        );
      })
      .catch((err) => console.error(err));
  };

  useEffect(() => {
    setHabit(habitList[0]);
  }, [habitList]);

  // runs on start and when habit updates
  useEffect(() => {
    if (habit && habit.id) {
      getDates();
    }
  }, [habit]);

  useEffect(() => {
    const myHabit = habitList.find((habit) => habit.name === url);
    if (myHabit !== undefined) {
      setHabit(myHabit);
    }
  }, [url, habitList]);

  return (
    <>
      {habit && (
        <div className="habit">
          {addOpen && (
            <AddActivity
              habit_id={habit.id}
              unit={habit.unit}
              setAddOpen={setAddOpen}
              getDates={getDates}
            />
          )}
          {editOpen && (
            <EditHabit
              habit={habit}
              setHabit={setHabit}
              setEditOpen={setEditOpen}
              updateHabits={updateHabits}
              habitList={habitList}
            />
          )}
          <div className="habit__header">
            <h1 className="habit__title">{habit.name}</h1>
            <img
              className="habit__icon"
              src={edit}
              alt="edit habit"
              onClick={() => {
                setEditOpen(true);
              }}
            />
            <img
              className="habit__icon"
              src={add}
              alt="add activity"
              onClick={() => {
                setAddOpen(true);
              }}
            />
          </div>
          <div className="habit__about">
            <div className="why">
              <h3 className="why__title">why</h3>
              <p className="why__text">{habit.why}</p>
            </div>
            <div className="goal">
              <h3 className="goal__title">goal:</h3>
              <div className="goal__list">
                <p className="goal__choice">daily</p>
                <p className="goal__choice">weekly</p>
                <p className="goal__choice goal__choice--active">
                  occasionally
                </p>
              </div>
            </div>
          </div>
          {dates && (
            <CalendarHeatmap
              startDate={shiftDate(today, -(DAYS_SHOWN - 1))}
              endDate={today}
              values={dates}
              classForValue={(value) => {
                if (!value) {
                  return "color-empty";
                }
                return `color-github-${countToLevel(value.count, maxCount)}`;
              }}
              tooltipDataAttrs={(value) => {
                if (!value || !value.date) {
                  return {};
                }
                return {
                  "data-tip": `${value.date}: ${value.count} ${habit.unit}`,
                };
              }}
              showWeekdayLabels={true}
            />
          )}
          <ReactTooltip />
          <div className="pomodoro-box">
            <Pomodoro />
          </div>
        </div>
      )}
    </>
  );
};

/* ********
 * helpers
 * ********/

function getRange(count) {
  return Array.from({ length: count }, (_, i) => i);
}

export default Habit;
