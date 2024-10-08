import { useCallback, useMemo, useState, useEffect } from "react";
import dayjs from "dayjs";
import axios from "axios";
// import { createMockData } from "./mock/appMock";
import { ParsedDatesRange } from "./utils/getDatesRange";
import { ConfigFormValues, SchedulerProjectData } from "./types/global";
// import ConfigPanel from "./components/ConfigPanel";
import { StyledSchedulerFrame } from "./styles";
import { Scheduler } from ".";

function App() {
  const [values, setValues] = useState<ConfigFormValues>({
    peopleCount: 15,
    projectsPerYear: 5,
    yearsCovered: 0,
    startDate: undefined,
    maxRecordsPerPage: 50,
    isFullscreen: true
  });

  const { peopleCount, projectsPerYear, yearsCovered, isFullscreen, maxRecordsPerPage } = values;

  // const mocked = useMemo(
  //   () => createMockData(+peopleCount, +yearsCovered, +projectsPerYear),
  //   [peopleCount, projectsPerYear, yearsCovered]
  // );

  const [mocked, setMocked] = useState<SchedulerProjectData[]>([]);

  useEffect(() => {
    // Fetch data from the API
    const fetchData = async () => {
      try {
        const configResponse = await axios.get('./config.json');
        const apiBaseUrl = configResponse.data.apiBaseUrl;
        const response = await axios.get(`${apiBaseUrl}/projects/doing`);
        setMocked(response.data); // Assuming response.data is an array of SchedulerProjectData
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    };

    fetchData();
  }, []);

  //console.log("Mocked Data: ", mocked);

  const [range, setRange] = useState<ParsedDatesRange>({
    startDate: new Date(),
    endDate: new Date()
  });

  const handleRangeChange = useCallback((range: ParsedDatesRange) => {
    setRange(range);
  }, []);

  const filteredData = useMemo(
    () =>
      mocked.map((project) => ({
        ...project,
        data: project.data.filter(
          (task) =>
            dayjs(task.startDate).isBetween(range.startDate, range.endDate) ||
            dayjs(task.endDate).isBetween(range.startDate, range.endDate) ||
            (dayjs(task.startDate).isBefore(range.startDate, "day") &&
              dayjs(task.endDate).isAfter(range.endDate, "day"))
        )
      })),
    [mocked, range.endDate, range.startDate]
  );

  const handleFilterData = () => console.log(`Filters button was clicked.`);

  const handleTileClick = (data: SchedulerProjectData) =>
    console.log(
      `Item ${data.title} - ${data.subtitle} was clicked. \n==============\nStart date: ${data.startDate} \n==============\nEnd date: ${data.endDate}\n==============\nOccupancy: ${data.occupancy}`
    );

  return (
    <>
      {/* <ConfigPanel values={values} onSubmit={setValues} /> */}
      {isFullscreen ? (
        <Scheduler
          startDate={values.startDate ? new Date(values.startDate).toISOString() : undefined}
          onRangeChange={handleRangeChange}
          data={filteredData}
          isLoading={false}
          onTileClick={handleTileClick}
          onFilterData={handleFilterData}
          config={{ zoom: 0, maxRecordsPerPage: maxRecordsPerPage }}
          onItemClick={(data) => console.log("clicked: ", data)}
        />
      ) : (
        <StyledSchedulerFrame>
          <Scheduler
            startDate={values.startDate ? new Date(values.startDate).toISOString() : undefined}
            onRangeChange={handleRangeChange}
            isLoading={false}
            data={filteredData}
            onTileClick={handleTileClick}
            onFilterData={handleFilterData}
            onItemClick={(data) => console.log("clicked: ", data)}
          />
        </StyledSchedulerFrame>
      )}
    </>
  );
}

export default App;
