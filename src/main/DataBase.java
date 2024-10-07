package main;

import Interfaces.IDataBase;
import Interfaces.IWork;
import java.util.HashMap;
import java.util.Map;

public class DataBase implements IDataBase<String, IWork<String, String>> {
    private Map<String , IWork<String, String>> workMap;

    public DataBase() {
        this.workMap = new HashMap<String, IWork<String, String>>();
    }

    @Override
    public void add(String key, IWork<String, String> value) {
        this.workMap.put(key, value);
    }

    @Override
    public IWork<String, String> get(String key) {
        return this.workMap.get(key);
    }

    @Override
    public void remove(String key) {
        this.workMap.remove(key);
    }

    @Override
    public void put(String key, IWork<String, String> value) {
        this.workMap.replace(key, value);
    }
}