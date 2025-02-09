import React, { useEffect, useState } from "react";
import { Table, Spinner, Alert } from "react-bootstrap";
import { parseISO, format } from "date-fns";
import axios from "axios";

interface DataType {
  ip: string;  // Первичный ключ
  ping_time: number;  // Время пинга в мс
  last_successful: string;  // Дата последнего успешного пинга
}

const DataTable: React.FC = () => {
  const [data, setData] = useState<DataType[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Функция для загрузки данных
  const fetchData = () => {
    axios.get("http://localhost:8080/v1/containers/")
      .then(response => {
        setData(response.data);
        setError(null); // Сброс ошибки, если данные успешно загружены
      })
      .catch(error => {
        setError("Ошибка загрузки данных");
        console.error(error);
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    // Загружаем данные сразу при монтировании компонента
    fetchData();

    // Устанавливаем интервал для обновления данных каждые 10 секунд
    const intervalId = setInterval(fetchData, 10000);

    // Очистка интервала при размонтировании компонента
    return () => clearInterval(intervalId);
  }, []); // Пустой массив зависимостей, чтобы эффект выполнялся только при монтировании и размонтировании

  return (
    <div className="container mt-4">
      <h2>Статистика пинга</h2>
      {loading && <Spinner animation="border" />}
      {error && <Alert variant="danger">{error}</Alert>}
      {!loading && !error && (
        <Table striped bordered hover>
          <thead>
            <tr>
              <th>IP-адрес</th>
              <th>Время пинга (мс)</th>
              <th>Последний успешный пинг</th>
            </tr>
          </thead>
          <tbody>
            {data.map((item) => (
              <tr key={item.ip}>
                <td>{item.ip}</td>
                <td>{item.ping_time} мс</td>
                <td>
                  {item.last_successful
                    ? format(parseISO(item.last_successful), "dd.MM.yyyy HH:mm")
                    : "Нет данных"}
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      )}
    </div>
  );
};

export default DataTable;